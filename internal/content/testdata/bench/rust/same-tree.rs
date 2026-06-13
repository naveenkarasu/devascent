use std::rc::Rc;
use std::cell::RefCell;

fn is_same_tree(p: Option<Rc<RefCell<TreeNode>>>, q: Option<Rc<RefCell<TreeNode>>>) -> bool {
    match (p, q) {
        (None, None) => true,
        (Some(a), Some(b)) => {
            let an = a.borrow();
            let bn = b.borrow();
            if an.val != bn.val {
                return false;
            }
            is_same_tree(an.left.clone(), bn.left.clone())
                && is_same_tree(an.right.clone(), bn.right.clone())
        }
        _ => false,
    }
}
