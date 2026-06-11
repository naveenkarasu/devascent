function is_strong_password(s) {
  if (s.length < 8) return false;
  const hasUpper = /[A-Z]/.test(s);
  const hasLower = /[a-z]/.test(s);
  const hasDigit = /[0-9]/.test(s);
  const hasSpecial = /[^a-zA-Z0-9]/.test(s);
  return hasUpper && hasLower && hasDigit && hasSpecial;
}
