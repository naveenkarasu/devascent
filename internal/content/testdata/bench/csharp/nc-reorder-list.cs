public class Solution {
    public ListNode reorder_list(ListNode head) {
        if (head == null || head.next == null) return head;
        // Find middle
        var slow = head; var fast = head;
        while (fast.next != null && fast.next.next != null) {
            slow = slow.next;
            fast = fast.next.next;
        }
        // Reverse second half
        var second = slow.next;
        slow.next = null;
        ListNode prev = null;
        while (second != null) {
            var nxt = second.next;
            second.next = prev;
            prev = second;
            second = nxt;
        }
        // Merge
        var first = head;
        while (prev != null) {
            var n1 = first.next;
            var n2 = prev.next;
            first.next = prev;
            prev.next = n1;
            first = n1;
            prev = n2;
        }
        return head;
    }
}
