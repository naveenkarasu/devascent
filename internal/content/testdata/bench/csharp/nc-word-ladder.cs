using System.Collections.Generic;

public class Solution {
    public long ladder_length(string begin_word, string end_word, string[] word_list) {
        var wordSet = new HashSet<string>(word_list);
        if (!wordSet.Contains(end_word)) return 0;
        var queue = new Queue<(string word, long length)>();
        queue.Enqueue((begin_word, 1));
        var visited = new HashSet<string>();
        visited.Add(begin_word);
        while (queue.Count > 0) {
            var (word, length) = queue.Dequeue();
            char[] chars = word.ToCharArray();
            for (int i = 0; i < chars.Length; i++) {
                char orig = chars[i];
                for (char c = 'a'; c <= 'z'; c++) {
                    chars[i] = c;
                    string candidate = new string(chars);
                    if (candidate == end_word) return length + 1;
                    if (wordSet.Contains(candidate) && !visited.Contains(candidate)) {
                        visited.Add(candidate);
                        queue.Enqueue((candidate, length + 1));
                    }
                }
                chars[i] = orig;
            }
        }
        return 0;
    }
}
