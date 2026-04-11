import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/api/stream_helpers.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'settings_transport_provider.dart';

/// Parameters for fetching a single setting.
class SettingKey {
  const SettingKey({
    required this.name,
    required this.object,
    required this.objectId,
    this.lang = '',
    this.module = '',
  });

  final String name;
  final String object;
  final String objectId;
  final String lang;
  final String module;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SettingKey &&
          name == other.name &&
          object == other.object &&
          objectId == other.objectId &&
          lang == other.lang &&
          module == other.module;

  @override
  int get hashCode => Object.hash(name, object, objectId, lang, module);
}

/// Parameters for listing settings.
class SettingListParams {
  const SettingListParams({
    required this.object,
    required this.objectId,
    this.module = '',
  });

  final String object;
  final String objectId;
  final String module;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SettingListParams &&
          object == other.object &&
          objectId == other.objectId &&
          module == other.module;

  @override
  int get hashCode => Object.hash(object, objectId, module);
}

/// Get a single setting by key.
final settingByKeyProvider =
    FutureProvider.family<SettingObject, SettingKey>((ref, key) async {
  final client = ref.watch(settingsServiceClientProvider);
  final request = GetRequest()
    ..setting = (Setting()
      ..name = key.name
      ..object = key.object
      ..objectId = key.objectId
      ..lang = key.lang
      ..module = key.module);
  final response = await client.get(request);
  return response.data;
});

/// List settings for a given object/instance, optionally filtered by module.
final settingsListProvider = FutureProvider.family<List<SettingObject>,
    SettingListParams>((ref, params) async {
  final client = ref.watch(settingsServiceClientProvider);
  final request = ListRequest()
    ..setting = (Setting()
      ..object = params.object
      ..objectId = params.objectId
      ..module = params.module);
  final stream = client.list(request);
  return collectStream<ListResponse, SettingObject>(
    stream,
    extract: (r) => r.data,
  );
});

/// Search settings by query string.
final settingsSearchProvider =
    FutureProvider.family<List<SettingObject>, String>((ref, query) async {
  final client = ref.watch(settingsServiceClientProvider);
  final request = ListRequest()
    ..setting = (Setting()..name = query);
  final stream = client.search(request);
  return collectStream<SearchResponse, SettingObject>(
    stream,
    extract: (r) => r.data,
  );
});

/// Notifier for setting mutations (set / update).
class SettingsNotifier extends StateNotifier<AsyncValue<void>> {
  SettingsNotifier(this._client) : super(const AsyncValue.data(null));
  final SettingsServiceClient _client;

  Future<SettingObject> set(SetRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.set(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final settingsNotifierProvider =
    StateNotifierProvider<SettingsNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(settingsServiceClientProvider);
  return SettingsNotifier(client);
});
