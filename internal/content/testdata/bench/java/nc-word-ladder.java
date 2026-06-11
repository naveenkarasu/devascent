import java.util.*;

class Solution {
    public long ladder_length(String begin_word, String end_word, String[] word_list) {
        Set<String> wordSet = new HashSet<>(Arrays.asList(word_list));
        if (!wordSet.contains(end_word)) return 0;
        Queue<String[]> queue = new LinkedList<>();
        queue.add(new String[]{begin_word, "1"});
        Set<String> visited = new HashSet<>();
        visited.add(begin_word);
        while (!queue.isEmpty()) {
            String[] cur = queue.poll();
            String word = cur[0];
            long length = Long.parseLong(cur[1]);
            char[] chars = word.toCharArray();
            for (int i = 0; i < chars.length; i++) {
                char orig = chars[i];
                for (char c = 'a'; c <= 'z'; c++) {
                    chars[i] = c;
                    String candidate = new String(chars);
                    if (candidate.equals(end_word)) return length + 1;
                    if (wordSet.contains(candidate) && !visited.contains(candidate)) {
                        visited.add(candidate);
                        queue.add(new String[]{candidate, String.valueOf(length + 1)});
                    }
                }
                chars[i] = orig;
            }
        }
        return 0;
    }
}
