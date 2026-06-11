import java.util.*;

class Solution {
    public Object[] detect_squares_ops(Object[][] operations) {
        Map<String, Long> cnt = new HashMap<>(); // "x,y" -> count
        List<long[]> points = new ArrayList<>();
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            Object[] op = operations[i];
            String name = (String) op[0];
            Object[] pt = (Object[]) op[1];
            long px = ((Number) pt[0]).longValue();
            long py = ((Number) pt[1]).longValue();
            if (name.equals("add")) {
                String key = px + "," + py;
                if (!cnt.containsKey(key)) points.add(new long[]{px, py});
                cnt.merge(key, 1L, Long::sum);
                out[i] = null;
            } else { // count
                long total = 0;
                for (long[] p : points) {
                    long x = p[0], y = p[1];
                    if (Math.abs(x - px) == Math.abs(y - py) && x != px && y != py) {
                        long c = cnt.get(x + "," + y);
                        long c1 = cnt.getOrDefault(px + "," + y, 0L);
                        long c2 = cnt.getOrDefault(x + "," + py, 0L);
                        total += c * c1 * c2;
                    }
                }
                out[i] = total;
            }
        }
        return out;
    }
}
