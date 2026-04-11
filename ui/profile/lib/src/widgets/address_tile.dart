import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter/material.dart';

/// A card widget for displaying an [AddressObject].
class AddressTile extends StatelessWidget {
  const AddressTile({
    super.key,
    required this.address,
    this.onTap,
    this.trailing,
  });

  final AddressObject address;
  final VoidCallback? onTap;
  final Widget? trailing;

  String _formatAddress() {
    final parts = <String>[];
    if (address.house.isNotEmpty) parts.add(address.house);
    if (address.street.isNotEmpty) parts.add(address.street);
    if (address.area.isNotEmpty) parts.add(address.area);
    if (address.city.isNotEmpty) parts.add(address.city);
    if (address.country.isNotEmpty) parts.add(address.country);
    return parts.isNotEmpty ? parts.join(', ') : 'No address details';
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final hasCoords = address.latitude != 0 || address.longitude != 0;

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: theme.colorScheme.outlineVariant),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(
                  Icons.location_on_outlined,
                  size: 20,
                  color: theme.colorScheme.onPrimaryContainer,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (address.name.isNotEmpty)
                      Text(
                        address.name,
                        style: theme.textTheme.titleSmall?.copyWith(
                          fontWeight: FontWeight.w600,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    Text(
                      _formatAddress(),
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 6),
                    Row(
                      children: [
                        if (address.postcode.isNotEmpty)
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 8,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: theme.colorScheme.secondaryContainer,
                              borderRadius: BorderRadius.circular(6),
                            ),
                            child: Text(
                              address.postcode,
                              style: theme.textTheme.labelSmall?.copyWith(
                                color:
                                    theme.colorScheme.onSecondaryContainer,
                              ),
                            ),
                          ),
                        if (hasCoords) ...[
                          if (address.postcode.isNotEmpty)
                            const SizedBox(width: 8),
                          Icon(
                            Icons.gps_fixed,
                            size: 12,
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                          const SizedBox(width: 4),
                          Text(
                            '${address.latitude.toStringAsFixed(4)}, '
                            '${address.longitude.toStringAsFixed(4)}',
                            style: theme.textTheme.labelSmall?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                        ],
                      ],
                    ),
                  ],
                ),
              ),
              if (trailing != null) trailing!,
            ],
          ),
        ),
      ),
    );
  }
}
