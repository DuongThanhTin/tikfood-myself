const secretPatterns = [
  /\.env($|\.)/,
  /secrets\//,
  /credentials\//,
  /id_rsa/,
  /private[_-]?key/i,
  /\.(pem|key)$/
];

export function assertNotSecretPath(path: string): void {
  if (secretPatterns.some((pattern) => pattern.test(path))) {
    throw new Error(`Secret-like path blocked: ${path}`);
  }
}
