use std::rc::Rc;
use std::cell::RefCell;

fn is_subtree(root: Option<Rc<RefCell<TreeNode>>>, sub: Option<Rc<RefCell<TreeNode>>>) -> bool {
    fn same(a: &Option<Rc<RefCell<TreeNode>>>, b: &Option<Rc<RefCell<TreeNode>>>) -> bool {
        match (a, b) {
            (None, None) => true,
            (Some(x), Some(y)) => {
                let xb = x.borrow();
                let yb = y.borrow();
                if xb.val != yb.val {
                    return false;
                }
                same(&xb.left, &yb.left) && same(&xb.right, &yb.right)
            }
            _ => false,
        }
    }
    fn dfs(node: &Option<Rc<RefCell<TreeNode>>>, sub: &Option<Rc<RefCell<TreeNode>>>) -> bool {
        match node {
            None => false,
            Some(n) => {
                if same(node, sub) {
                    return true;
                }
                let nb = n.borrow();
                dfs(&nb.left, sub) || dfs(&nb.right, sub)
            }
        }
    }
    if sub.is_none() {
        return true;
    }
    dfs(&root, &sub)
}
