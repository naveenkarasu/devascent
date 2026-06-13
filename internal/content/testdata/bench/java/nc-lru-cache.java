import java.util.*;

class Solution {
    public Object[] lru_cache_ops(long capacity, Object[][] operations) {
        LinkedHashMap<Long, Long> d = new LinkedHashMap<>();
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            Object[] op = operations[i];
            String name = (String) op[0];
            if (name.equals("get")) {
                long key = ((Number) op[1]).longValue();
                if (!d.containsKey(key)) {
                    out[i] = -1L;
                } else {
                    long val = d.remove(key);
                    d.put(key, val);
                    out[i] = val;
                }
            } else { // put
                long key = ((Number) op[1]).longValue();
                long value = ((Number) op[2]).longValue();
                if (d.containsKey(key)) d.remove(key);
                d.put(key, value);
                if (d.size() > capacity) {
                    Iterator<Long> it = d.keySet().iterator();
                    it.next();
                    it.remove();
                }
                out[i] = null;
            }
        }
        return out;
    }
}
