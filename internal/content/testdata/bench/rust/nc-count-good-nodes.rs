use std::rc::Rc;
use std::cell::RefCell;

fn dfs(node: &Option<Rc<RefCell<TreeNode>>>, mx: i64) -> i64 {
    match node {
        None => 0,
        Some(n) => {
            let v = n.borrow().val as i64;
            let count = if v >= mx { 1 } else { 0 };
            let nmx = if v > mx { v } else { mx };
            count + dfs(&n.borrow().left, nmx) + dfs(&n.borrow().right, nmx)
        }
    }
}

fn good_nodes(root: Option<Rc<RefCell<TreeNode>>>) -> i64 {
    dfs(&root, i64::MIN)
}
