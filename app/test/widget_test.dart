import 'package:flutter_test/flutter_test.dart';

import 'package:chhootganit_app/main.dart';

void main() {
  test('evaluate: 5% off holds profit with +33.3% units', () {
    final r = evaluate(100, 20, 5);
    expect(r.feasible, true);
    expect(r.extraPct, closeTo(100 / 3, 1e-6));
  });

  test('evaluate: discount past margin is infeasible', () {
    expect(evaluate(100, 20, 25).feasible, false);
  });

  testWidgets('renders the counter card', (tester) async {
    await tester.pumpWidget(const ChhootganitApp());
    expect(find.text('Keep this at the counter'), findsOneWidget);
  });
}
