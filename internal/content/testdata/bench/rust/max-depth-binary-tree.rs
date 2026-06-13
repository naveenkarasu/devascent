use std::rc::Rc;
use std::cell::RefCell;

fn max_depth(root: Option<Rc<RefCell<TreeNode>>>) -> i64 {
    match root {
        None => 0,
        Some(node) => {
            let n = node.borrow();
            1 + max_depth(n.left.clone()).max(max_depth(n.right.clone()))
        }
    }
}
