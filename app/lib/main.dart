import 'package:flutter/material.dart';

void main() => runApp(const ChhootganitApp());

/// Chhootganit — discount maths. Given profit per unit and a proposed discount,
/// how many more units must you sell to keep total profit? Mirrors the Go engine
/// and prints a threshold reference card for the counter.
class ChhootganitApp extends StatelessWidget {
  const ChhootganitApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Chhootganit',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFFB2662E), useMaterial3: true),
        home: const HomePage(),
      );
}

class Result {
  final double discountPerUnit, newProfit, extraPct;
  final bool feasible;
  const Result(this.discountPerUnit, this.newProfit, this.extraPct, this.feasible);
}

/// evaluate mirrors backend/cost.go.
Result evaluate(double price, double profit, double discountPct) {
  final d = price * discountPct / 100;
  final np = profit - d;
  if (np <= 0) return Result(d, np, 0, false);
  return Result(d, np, (profit / np - 1) * 100, true);
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _price = TextEditingController(text: '100');
  final _profit = TextEditingController(text: '20');

  @override
  Widget build(BuildContext context) {
    final price = double.tryParse(_price.text.trim()) ?? 0;
    final profit = double.tryParse(_profit.text.trim()) ?? 0;
    const discounts = <double>[2, 5, 10, 15, 20, 25];
    return Scaffold(
      appBar: AppBar(
        title: const Text('Chhootganit · discount card'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: ListView(padding: const EdgeInsets.all(16), children: [
        Row(children: [
          Expanded(child: _f(_price, 'Selling price ₹')),
          const SizedBox(width: 12),
          Expanded(child: _f(_profit, 'Profit per unit ₹')),
        ]),
        const SizedBox(height: 16),
        const Text('Keep this at the counter', style: TextStyle(fontWeight: FontWeight.w600)),
        const SizedBox(height: 8),
        Card(
          color: Theme.of(context).colorScheme.primaryContainer,
          child: Column(children: [
            const ListTile(
              dense: true,
              title: Text('Discount', style: TextStyle(fontWeight: FontWeight.bold)),
              trailing: Text('Extra units to sell', style: TextStyle(fontWeight: FontWeight.bold)),
            ),
            for (final d in discounts) _rowFor(price, profit, d),
          ]),
        ),
        const SizedBox(height: 12),
        const Text('Grant a discount only if you can realistically sell that many more.',
            style: TextStyle(fontSize: 13)),
      ]),
    );
  }

  Widget _rowFor(double price, double profit, double d) {
    final r = evaluate(price, profit, d);
    return ListTile(
      dense: true,
      title: Text('${d.toStringAsFixed(0)}%'),
      trailing: Text(
        r.feasible ? '+${r.extraPct.toStringAsFixed(1)}%' : 'impossible',
        style: TextStyle(
          fontWeight: FontWeight.bold,
          color: r.feasible ? null : Colors.red,
        ),
      ),
    );
  }

  Widget _f(TextEditingController c, String label) => TextField(
        controller: c,
        keyboardType: const TextInputType.numberWithOptions(decimal: true),
        decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
        onChanged: (_) => setState(() {}),
      );
}
