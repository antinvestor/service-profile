import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/contact_providers.dart';

/// Dialog for OTP-based contact verification.
///
/// Two-step flow: 1) Send a verification code, 2) Enter the code to verify.
class ContactVerificationDialog extends ConsumerStatefulWidget {
  const ContactVerificationDialog({
    super.key,
    required this.profileId,
    required this.contactId,
  });

  final String profileId;
  final String contactId;

  @override
  ConsumerState<ContactVerificationDialog> createState() =>
      _ContactVerificationDialogState();
}

class _ContactVerificationDialogState
    extends ConsumerState<ContactVerificationDialog> {
  final _codeController = TextEditingController();
  String? _verificationId;
  bool _isLoading = false;
  String? _error;
  bool _verified = false;

  @override
  void dispose() {
    _codeController.dispose();
    super.dispose();
  }

  Future<void> _sendCode() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final notifier = ref.read(contactNotifierProvider.notifier);
      await notifier.createVerification(
        CreateContactVerificationRequest(
          id: widget.profileId,
          contactId: widget.contactId,
          durationToExpire: '5m',
        ),
      );
      setState(() {
        _verificationId = widget.contactId;
        _isLoading = false;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Verification code sent'),
            behavior: SnackBarBehavior.floating,
            width: 260,
          ),
        );
      }
    } catch (e) {
      setState(() {
        _error = friendlyError(e);
        _isLoading = false;
      });
    }
  }

  Future<void> _verify() async {
    final code = _codeController.text.trim();
    if (code.isEmpty) {
      setState(() => _error = 'Please enter the verification code.');
      return;
    }
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final notifier = ref.read(contactNotifierProvider.notifier);
      final passed = await notifier.checkVerification(
        CheckVerificationRequest(
          id: _verificationId ?? widget.contactId,
          code: code,
        ),
      );
      if (passed) {
        setState(() {
          _verified = true;
          _isLoading = false;
        });
      } else {
        setState(() {
          _error = 'Verification failed. Please check the code and try again.';
          _isLoading = false;
        });
      }
    } catch (e) {
      setState(() {
        _error = friendlyError(e);
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AlertDialog(
      title: Text(_verified ? 'Verified' : 'Verify Contact'),
      content: SizedBox(
        width: 340,
        child: _verified
            ? Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.check_circle, size: 64, color: Colors.green),
                  const SizedBox(height: 12),
                  Text(
                    'Contact verified successfully.',
                    style: theme.textTheme.bodyLarge,
                  ),
                ],
              )
            : Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Text(
                    'Enter the verification code sent to your contact.',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                  const SizedBox(height: 16),
                  TextField(
                    controller: _codeController,
                    decoration: const InputDecoration(
                      labelText: 'Verification Code',
                      hintText: 'Enter code...',
                    ),
                    keyboardType: TextInputType.number,
                    enabled: !_isLoading,
                  ),
                  if (_error != null) ...[
                    const SizedBox(height: 12),
                    Text(
                      _error!,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.error,
                      ),
                    ),
                  ],
                  if (_isLoading) ...[
                    const SizedBox(height: 16),
                    const Center(child: CircularProgressIndicator()),
                  ],
                ],
              ),
      ),
      actions: _verified
          ? [
              FilledButton(
                onPressed: () => Navigator.of(context).pop(true),
                child: const Text('Done'),
              ),
            ]
          : [
              TextButton(
                onPressed: _isLoading ? null : () => Navigator.of(context).pop(),
                child: const Text('Cancel'),
              ),
              OutlinedButton(
                onPressed: _isLoading ? null : _sendCode,
                child: const Text('Send Code'),
              ),
              FilledButton(
                onPressed: _isLoading ? null : _verify,
                child: const Text('Verify'),
              ),
            ],
    );
  }
}
