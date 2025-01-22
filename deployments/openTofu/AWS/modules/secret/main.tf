resource "aws_secretsmanager_secret" "aws_secret" {
  name = var.secret_name
}

resource "aws_secretsmanager_secret_version" "aws_secret_version" {
  secret_id     = aws_secretsmanager_secret.aws_secret.id
  secret_string = jsonencode(var.secret_value)
}