// Placeholder for settings provider tests.
//
// Integration tests require a running settings service or mock transport.
// Unit tests for SettingKey/SettingListParams equality and SettingsScope
// are straightforward to add here.

import 'package:flutter_test/flutter_test.dart';
import 'package:antinvestor_ui_settings/antinvestor_ui_settings.dart';

void main() {
  group('SettingKey equality', () {
    test('equal keys match', () {
      const a = SettingKey(
        name: 'theme',
        object: 'app',
        objectId: '1',
        lang: 'en',
        module: 'ui',
      );
      const b = SettingKey(
        name: 'theme',
        object: 'app',
        objectId: '1',
        lang: 'en',
        module: 'ui',
      );
      expect(a, equals(b));
      expect(a.hashCode, equals(b.hashCode));
    });

    test('different keys do not match', () {
      const a = SettingKey(name: 'theme', object: 'app', objectId: '1');
      const b = SettingKey(name: 'locale', object: 'app', objectId: '1');
      expect(a, isNot(equals(b)));
    });
  });

  group('SettingListParams equality', () {
    test('equal params match', () {
      const a = SettingListParams(object: 'tenant', objectId: '42');
      const b = SettingListParams(object: 'tenant', objectId: '42');
      expect(a, equals(b));
      expect(a.hashCode, equals(b.hashCode));
    });
  });

  group('SettingsScope', () {
    test('copyWith returns updated scope', () {
      const scope = SettingsScope(object: 'app', objectId: '1', lang: 'en');
      final updated = scope.copyWith(lang: 'fr');
      expect(updated.lang, 'fr');
      expect(updated.object, 'app');
      expect(updated.objectId, '1');
    });

    test('equality', () {
      const a = SettingsScope(object: 'app', objectId: '1', lang: 'en');
      const b = SettingsScope(object: 'app', objectId: '1', lang: 'en');
      expect(a, equals(b));
    });
  });
}
