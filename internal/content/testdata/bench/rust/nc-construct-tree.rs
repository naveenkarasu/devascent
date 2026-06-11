use std::rc::Rc;
use std::cell::RefCell;
use std::collections::HashMap;

fn build_tree(preorder: Vec<i64>, inorder: Vec<i64>) -> Option<Rc<RefCell<TreeNode>>> {
    let mut idx: HashMap<i64, i64> = HashMap::new();
    for (i, &v) in inorder.iter().enumerate() { idx.insert(v, i as i64); }
    fn build(preorder: &Vec<i64>, idx: &HashMap<i64, i64>, pre: &mut usize, lo: i64, hi: i64) -> Option<Rc<RefCell<TreeNode>>> {
        if lo > hi { return None; }
        let val = preorder[*pre];
        *pre += 1;
        let node = Rc::new(RefCell::new(TreeNode::new(val as i32)));
        let mid = idx[&val];
        let l = build(preorder, idx, pre, lo, mid - 1);
        let r = build(preorder, idx, pre, mid + 1, hi);
        node.borrow_mut().left = l;
        node.borrow_mut().right = r;
        Some(node)
    }
    let mut pre = 0usize;
    let n = inorder.len() as i64;
    build(&preorder, &idx, &mut pre, 0, n - 1)
}
