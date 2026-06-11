use std::rc::Rc;
use std::cell::RefCell;

fn is_balanced(root: Option<Rc<RefCell<TreeNode>>>) -> bool {
    fn height(node: &Option<Rc<RefCell<TreeNode>>>) -> i64 {
        match node {
            None => 0,
            Some(n) => {
                let nb = n.borrow();
                let l = height(&nb.left);
                if l < 0 {
                    return -1;
                }
                let r = height(&nb.right);
                if r < 0 {
                    return -1;
                }
                if (l - r).abs() > 1 {
                    return -1;
                }
                1 + l.max(r)
            }
        }
    }
    height(&root) >= 0
}
