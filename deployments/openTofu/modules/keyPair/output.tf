output "aws_key_pair_name" {
  value = aws_key_pair.generated_key.key_name
}
output "aws_key_pem" {
  value = tls_private_key.key.private_key_pem
}
