import java.util.*;

class Solution {
    public Object[] time_map_ops(Object[][] operations) {
        Map<String, List<long[]>> times = new HashMap<>();   // key -> list of timestamps (with index into vals)
        Map<String, List<String>> vals = new HashMap<>();
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            Object[] op = operations[i];
            String name = (String) op[0];
            if (name.equals("set")) {
                String key = (String) op[1];
                String value = (String) op[2];
                long ts = ((Number) op[3]).longValue();
                times.computeIfAbsent(key, k -> new ArrayList<>()).add(new long[]{ts});
                vals.computeIfAbsent(key, k -> new ArrayList<>()).add(value);
                out[i] = null;
            } else { // get
                String key = (String) op[1];
                long ts = ((Number) op[2]).longValue();
                List<long[]> ts_list = times.get(key);
                if (ts_list == null) {
                    out[i] = "";
                } else {
                    // bisect_right: largest index with timestamp <= ts
                    int lo = 0, hi = ts_list.size();
                    while (lo < hi) {
                        int mid = (lo + hi) >>> 1;
                        if (ts_list.get(mid)[0] <= ts) lo = mid + 1; else hi = mid;
                    }
                    out[i] = (lo > 0) ? vals.get(key).get(lo - 1) : "";
                }
            }
        }
        return out;
    }
}
