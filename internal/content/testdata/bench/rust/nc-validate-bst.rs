use std::rc::Rc;
use std::cell::RefCell;

fn is_valid_bst(root: Option<Rc<RefCell<TreeNode>>>) -> bool {
    fn valid(node: &Option<Rc<RefCell<TreeNode>>>, lo: i64, hi: i64) -> bool {
        match node {
            None => true,
            Some(n) => {
                let nb = n.borrow();
                let v = nb.val as i64;
                if !(lo < v && v < hi) {
                    return false;
                }
                valid(&nb.left, lo, v) && valid(&nb.right, v, hi)
            }
        }
    }
    valid(&root, i64::MIN, i64::MAX)
}
