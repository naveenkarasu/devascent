use std::rc::Rc;
use std::cell::RefCell;

fn codec_roundtrip(root: Option<Rc<RefCell<TreeNode>>>) -> Option<Rc<RefCell<TreeNode>>> {
    fn ser(node: &Option<Rc<RefCell<TreeNode>>>, out: &mut Vec<String>) {
        match node {
            None => out.push("#".to_string()),
            Some(n) => {
                let (v, l, r) = { let b = n.borrow(); (b.val, b.left.clone(), b.right.clone()) };
                out.push(v.to_string());
                ser(&l, out);
                ser(&r, out);
            }
        }
    }
    fn de(vals: &[String], i: &mut usize) -> Option<Rc<RefCell<TreeNode>>> {
        let v = vals[*i].clone();
        *i += 1;
        if v == "#" { return None; }
        let node = Rc::new(RefCell::new(TreeNode::new(v.parse::<i32>().unwrap())));
        let l = de(vals, i);
        let r = de(vals, i);
        node.borrow_mut().left = l;
        node.borrow_mut().right = r;
        Some(node)
    }
    let mut out = Vec::new();
    ser(&root, &mut out);
    let data = out.join(",");
    let vals: Vec<String> = data.split(',').map(|s| s.to_string()).collect();
    let mut i = 0;
    de(&vals, &mut i)
}
