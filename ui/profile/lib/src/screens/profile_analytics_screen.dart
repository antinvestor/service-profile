import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/profile_providers.dart';

/// Profile service dashboard with KPI cards showing profile statistics.
class ProfileAnalyticsScreen extends ConsumerWidget {
  const ProfileAnalyticsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final asyncProfiles = ref.watch(profileSearchProvider(''));

    final profiles =
        asyncProfiles.whenOrNull(data: (d) => d) ?? <ProfileObject>[];
    final totalCount = profiles.length;
    final personCount =
        profiles.where((p) => p.type == ProfileType.PERSON).length;
    final institutionCount =
        profiles.where((p) => p.type == ProfileType.INSTITUTION).length;
    final botCount =
        profiles.where((p) => p.type == ProfileType.BOT).length;
    final activeCount =
        profiles.where((p) => p.state == STATE.ACTIVE).length;
    final totalContacts =
        profiles.fold<int>(0, (sum, p) => sum + p.contacts.length);
    final verifiedContacts = profiles.fold<int>(
        0,
        (sum, p) =>
            sum + p.contacts.where((c) => c.verified).length);

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          Row(
            children: [
              Icon(Icons.analytics_outlined,
                  size: 28, color: theme.colorScheme.primary),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('Profile Service',
                        style: theme.textTheme.headlineSmall
                            ?.copyWith(fontWeight: FontWeight.w600)),
                    Text('Service analytics and overview',
                        style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant)),
                  ],
                ),
              ),
              if (asyncProfiles.isLoading)
                const SizedBox(
                  width: 24,
                  height: 24,
                  child: CircularProgressIndicator(strokeWidth: 2),
                ),
            ],
          ),
          const SizedBox(height: 24),

          // KPI cards grid
          LayoutBuilder(
            builder: (context, constraints) {
              final crossAxisCount = constraints.maxWidth > 900
                  ? 4
                  : constraints.maxWidth > 600
                      ? 3
                      : 2;
              return Wrap(
                spacing: 16,
                runSpacing: 16,
                children: [
                  _KpiCard(
                    label: 'Total Profiles',
                    value: '$totalCount',
                    icon: Icons.people_outlined,
                    color: theme.colorScheme.primary,
                    width: (constraints.maxWidth -
                            (crossAxisCount - 1) * 16) /
                        crossAxisCount,
                  ),
                  _KpiCard(
                    label: 'Person Profiles',
                    value: '$personCount',
                    icon: Icons.person_outlined,
                    color: Colors.blue,
                    width: (constraints.maxWidth -
                            (crossAxisCount - 1) * 16) /
                        crossAxisCount,
                  ),
                  _KpiCard(
                    label: 'Institution Profiles',
                    value: '$institutionCount',
                    icon: Icons.business_outlined,
                    color: Colors.green,
                    width: (constraints.maxWidth -
                            (crossAxisCount - 1) * 16) /
                        crossAxisCount,
                  ),
                  _KpiCard(
                    label: 'Bot Profiles',
                    value: '$botCount',
                    icon: Icons.smart_toy_outlined,
                    color: Colors.orange,
                    width: (constraints.maxWidth -
                            (crossAxisCount - 1) * 16) /
                        crossAxisCount,
                  ),
                  _KpiCard(
                    label: 'Active Profiles',
                    value: '$activeCount',
                    icon: Icons.check_circle_outline,
                    color: Colors.teal,
                    width: (constraints.maxWidth -
                            (crossAxisCount - 1) * 16) /
                        crossAxisCount,
                  ),
                  _KpiCard(
                    label: 'Total Contacts',
                    value: '$totalContacts',
                    icon: Icons.contact_phone_outlined,
                    color: Colors.indigo,
                    width: (constraints.maxWidth -
                            (crossAxisCount - 1) * 16) /
                        crossAxisCount,
                  ),
                  _KpiCard(
                    label: 'Verified Contacts',
                    value: '$verifiedContacts',
                    icon: Icons.verified_outlined,
                    color: Colors.purple,
                    width: (constraints.maxWidth -
                            (crossAxisCount - 1) * 16) /
                        crossAxisCount,
                  ),
                ],
              );
            },
          ),
          const SizedBox(height: 32),

          // Recent events (placeholder)
          Text('Recent Events',
              style: theme.textTheme.titleMedium
                  ?.copyWith(fontWeight: FontWeight.w600)),
          const SizedBox(height: 12),
          Container(
            width: double.infinity,
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(12),
              border:
                  Border.all(color: theme.colorScheme.outlineVariant),
            ),
            child: Column(
              children: [
                _EventTile(
                  title: 'New profile registered',
                  timeAgo: '5 mins ago',
                  icon: Icons.person_add_outlined,
                  color: Colors.green,
                ),
                const Divider(height: 1),
                _EventTile(
                  title: 'Contact verification completed',
                  timeAgo: '20 mins ago',
                  icon: Icons.verified_outlined,
                  color: Colors.blue,
                ),
                const Divider(height: 1),
                _EventTile(
                  title: 'Profile merge executed',
                  timeAgo: '1 hour ago',
                  icon: Icons.merge_outlined,
                  color: Colors.orange,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _KpiCard extends StatelessWidget {
  const _KpiCard({
    required this.label,
    required this.value,
    required this.icon,
    required this.color,
    required this.width,
  });

  final String label;
  final String value;
  final IconData icon;
  final Color color;
  final double width;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return SizedBox(
      width: width,
      child: Card(
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: BorderSide(color: theme.colorScheme.outlineVariant),
        ),
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: color.withAlpha(25),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Icon(icon, size: 20, color: color),
                  ),
                  const Spacer(),
                ],
              ),
              const SizedBox(height: 16),
              Text(
                value,
                style: theme.textTheme.headlineMedium?.copyWith(
                  fontWeight: FontWeight.w700,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                label,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _EventTile extends StatelessWidget {
  const _EventTile({
    required this.title,
    required this.timeAgo,
    required this.icon,
    required this.color,
  });

  final String title;
  final String timeAgo;
  final IconData icon;
  final Color color;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ListTile(
      leading: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: color.withAlpha(25),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Icon(icon, size: 18, color: color),
      ),
      title: Text(title,
          style: theme.textTheme.bodyMedium
              ?.copyWith(fontWeight: FontWeight.w500)),
      trailing: Text(timeAgo,
          style: theme.textTheme.bodySmall
              ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
    );
  }
}
