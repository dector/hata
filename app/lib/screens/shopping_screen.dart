import 'package:flutter/material.dart';

import '../api/api_error.dart';
import '../api/hata_api_client.dart';
import '../models/shopping_item.dart';
import '../session/session.dart';
import '../theme/hata_colors.dart';
import '../widgets/shopping_tile.dart';

class ShoppingScreen extends StatefulWidget {
  const ShoppingScreen({
    super.key,
    required this.session,
    required this.apiClient,
    required this.onSessionInvalid,
    this.houseId,
  });

  final Session session;
  final HataApiClient apiClient;
  final Future<void> Function() onSessionInvalid;
  final String? houseId;

  @override
  State<ShoppingScreen> createState() => _ShoppingScreenState();
}

class _ShoppingScreenState extends State<ShoppingScreen> {
  static const _defaultListId = 'default';

  bool _isLoading = true;
  bool _isSaving = false;
  String? _errorMessage;
  String? _inlineErrorMessage;
  String? _houseId;
  List<ShoppingItem> _items = const [];
  final Set<String> _pendingItemIds = <String>{};

  @override
  void initState() {
    super.initState();
    _houseId = widget.houseId;
    _loadItems();
  }

  Future<void> _loadItems({bool showLoadingState = true}) async {
    if (showLoadingState) {
      setState(() {
        _isLoading = true;
        _errorMessage = null;
        _inlineErrorMessage = null;
      });
    } else {
      setState(() => _inlineErrorMessage = null);
    }

    try {
      final houseId = await _resolveHouseId();
      final items = await widget.apiClient.fetchShoppingItems(
        widget.session.serverUrl,
        widget.session.token,
        houseId,
        _defaultListId,
      );
      if (!mounted) return;
      setState(() {
        _houseId = houseId;
        _items = _sortItems(items);
        _isLoading = false;
        _errorMessage = null;
        _inlineErrorMessage = null;
      });
    } catch (error) {
      if (_isUnauthorized(error)) {
        await widget.onSessionInvalid();
        return;
      }
      if (!mounted) return;
      setState(() {
        _isLoading = false;
        if (_items.isEmpty) {
          _errorMessage = _messageFor(error);
        } else {
          _inlineErrorMessage = _messageFor(error);
        }
      });
    }
  }

  Future<String> _resolveHouseId() async {
    if (_houseId != null && _houseId!.isNotEmpty) return _houseId!;

    final houses = await widget.apiClient.fetchHouses(
      widget.session.serverUrl,
      widget.session.token,
    );
    final houseId = houses.firstOrNull?.id;
    if (houseId == null || houseId.isEmpty) {
      throw const ApiError('No house is available for shopping');
    }
    return houseId;
  }

  Future<void> _toggleItem(ShoppingItem item, bool checked) async {
    if (_pendingItemIds.contains(item.id)) return;
    final houseId = _houseId;
    if (houseId == null) return;

    setState(() => _pendingItemIds.add(item.id));
    try {
      final updated = await widget.apiClient.setShoppingItemChecked(
        widget.session.serverUrl,
        widget.session.token,
        houseId,
        _defaultListId,
        item.id,
        checked,
      );
      if (!mounted) return;
      setState(() {
        _pendingItemIds.remove(item.id);
        _replaceItem(updated);
      });
      _showToggleSnackBar(updated, checked);
    } catch (error) {
      if (_isUnauthorized(error)) {
        await widget.onSessionInvalid();
        return;
      }
      if (!mounted) return;
      setState(() => _pendingItemIds.remove(item.id));
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text(_messageFor(error))));
    }
  }

  Future<void> _undoToggle(ShoppingItem item, bool checked) async {
    final houseId = _houseId;
    if (houseId == null || _pendingItemIds.contains(item.id)) return;

    setState(() => _pendingItemIds.add(item.id));
    try {
      final updated = await widget.apiClient.setShoppingItemChecked(
        widget.session.serverUrl,
        widget.session.token,
        houseId,
        _defaultListId,
        item.id,
        !checked,
      );
      if (!mounted) return;
      setState(() {
        _pendingItemIds.remove(item.id);
        _replaceItem(updated);
      });
    } catch (error) {
      if (_isUnauthorized(error)) {
        await widget.onSessionInvalid();
        return;
      }
      if (!mounted) return;
      setState(() => _pendingItemIds.remove(item.id));
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text(_messageFor(error))));
    }
  }

  void _showToggleSnackBar(ShoppingItem item, bool checked) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('${item.name}: ${checked ? 'checked' : 'unchecked'}'),
        action: SnackBarAction(
          label: 'UNDO',
          onPressed: () => _undoToggle(item, checked),
        ),
      ),
    );
  }

  Future<void> _showAddItemDialog() async {
    final houseId = _houseId;
    if (houseId == null) return;

    final controller = TextEditingController();
    String? dialogError;
    await showDialog<void>(
      context: context,
      barrierDismissible: !_isSaving,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setDialogState) {
            final canSave = controller.text.trim().isNotEmpty && !_isSaving;
            return AlertDialog(
              title: const Text('Add item'),
              content: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TextField(
                    controller: controller,
                    autofocus: true,
                    textCapitalization: TextCapitalization.sentences,
                    decoration: InputDecoration(
                      labelText: 'Item name',
                      errorText: dialogError,
                    ),
                    onChanged: (_) => setDialogState(() {}),
                    onSubmitted: canSave
                        ? (_) => _createItem(
                            houseId,
                            controller.text,
                            setDialogState,
                            (message) => dialogError = message,
                          )
                        : null,
                  ),
                ],
              ),
              actions: [
                TextButton(
                  onPressed: _isSaving ? null : () => Navigator.pop(context),
                  child: const Text('Cancel'),
                ),
                FilledButton(
                  onPressed: canSave
                      ? () => _createItem(
                          houseId,
                          controller.text,
                          setDialogState,
                          (message) => dialogError = message,
                        )
                      : null,
                  child: _isSaving
                      ? const SizedBox.square(
                          dimension: 18,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('Save'),
                ),
              ],
            );
          },
        );
      },
    );
    controller.dispose();
  }

  Future<void> _createItem(
    String houseId,
    String rawName,
    StateSetter setDialogState,
    void Function(String?) setDialogError,
  ) async {
    final name = rawName.trim();
    if (name.isEmpty || _isSaving) return;

    setDialogState(() {
      _isSaving = true;
      setDialogError(null);
    });
    try {
      final item = await widget.apiClient.createShoppingItem(
        widget.session.serverUrl,
        widget.session.token,
        houseId,
        _defaultListId,
        name,
      );
      if (!mounted) return;
      setState(() {
        _items = _sortItems([..._items, item]);
        _isSaving = false;
      });
      Navigator.pop(context);
    } catch (error) {
      if (_isUnauthorized(error)) {
        await widget.onSessionInvalid();
        return;
      }
      if (!mounted) return;
      setDialogState(() {
        _isSaving = false;
        setDialogError(_messageFor(error));
      });
    }
  }

  void _replaceItem(ShoppingItem updated) {
    final nextItems = List<ShoppingItem>.from(_items);
    final index = nextItems.indexWhere((item) => item.id == updated.id);
    if (index == -1) {
      nextItems.add(updated);
    } else {
      nextItems[index] = updated;
    }
    _items = _sortItems(nextItems);
  }

  List<ShoppingItem> _sortItems(List<ShoppingItem> items) {
    return List<ShoppingItem>.from(items)..sort((a, b) {
      final checkedComparison = a.isChecked == b.isChecked
          ? 0
          : a.isChecked
          ? 1
          : -1;
      if (checkedComparison != 0) return checkedComparison;

      final nameComparison = a.name.toLowerCase().compareTo(
        b.name.toLowerCase(),
      );
      if (nameComparison != 0) return nameComparison;

      return a.uid.compareTo(b.uid);
    });
  }

  bool _isUnauthorized(Object error) =>
      error is ApiError && (error.statusCode == 401 || error.statusCode == 403);

  String _messageFor(Object error) =>
      error is ApiError ? error.message : 'Could not reach the server';

  @override
  Widget build(BuildContext context) {
    final hasLoaded = !_isLoading && _errorMessage == null;
    return Scaffold(
      appBar: AppBar(
        backgroundColor: HataColors.background,
        foregroundColor: Colors.white,
        title: const Text('Shopping'),
      ),
      floatingActionButton: hasLoaded
          ? FloatingActionButton(
              onPressed: _showAddItemDialog,
              child: const Icon(Icons.add),
            )
          : null,
      body: SafeArea(
        child: RefreshIndicator(
          color: HataColors.primary,
          backgroundColor: HataColors.surface,
          onRefresh: () => _loadItems(showLoadingState: false),
          child: _buildBody(),
        ),
      ),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const CustomScrollView(
        physics: AlwaysScrollableScrollPhysics(),
        slivers: [
          SliverFillRemaining(
            hasScrollBody: false,
            child: Center(child: CircularProgressIndicator()),
          ),
        ],
      );
    }

    if (_errorMessage != null) {
      return CustomScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          SliverFillRemaining(
            hasScrollBody: false,
            child: _MessageState(
              icon: Icons.cloud_off,
              title: 'Could not load shopping',
              message: _errorMessage!,
              actionLabel: 'Retry',
              onAction: () => _loadItems(),
            ),
          ),
        ],
      );
    }

    if (_items.isEmpty) {
      return const CustomScrollView(
        physics: AlwaysScrollableScrollPhysics(),
        slivers: [
          SliverFillRemaining(
            hasScrollBody: false,
            child: _MessageState(
              icon: Icons.shopping_basket_outlined,
              title: 'Nothing to buy',
              message: 'Your default shopping list is empty',
            ),
          ),
        ],
      );
    }

    return ListView.separated(
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 96),
      itemBuilder: (context, index) {
        if (index == 0 && _inlineErrorMessage != null) {
          return Text(
            _inlineErrorMessage!,
            style: const TextStyle(color: HataColors.errorVariant),
          );
        }
        final itemIndex = _inlineErrorMessage == null ? index : index - 1;
        final item = _items[itemIndex];
        return ShoppingTile(
          item: item,
          enabled: !_pendingItemIds.contains(item.id),
          onChanged: (value) => _toggleItem(item, value ?? false),
        );
      },
      separatorBuilder: (context, index) => const SizedBox(height: 12),
      itemCount: _items.length + (_inlineErrorMessage == null ? 0 : 1),
    );
  }
}

class _MessageState extends StatelessWidget {
  const _MessageState({
    required this.icon,
    required this.title,
    required this.message,
    this.actionLabel,
    this.onAction,
  });

  final IconData icon;
  final String title;
  final String message;
  final String? actionLabel;
  final VoidCallback? onAction;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: 48, color: HataColors.onSurfaceVariant),
          const SizedBox(height: 16),
          Text(
            title,
            textAlign: TextAlign.center,
            style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w700),
          ),
          const SizedBox(height: 8),
          Text(
            message,
            textAlign: TextAlign.center,
            style: const TextStyle(color: HataColors.onSurfaceVariant),
          ),
          if (actionLabel != null && onAction != null) ...[
            const SizedBox(height: 20),
            FilledButton(onPressed: onAction, child: Text(actionLabel!)),
          ],
        ],
      ),
    );
  }
}
