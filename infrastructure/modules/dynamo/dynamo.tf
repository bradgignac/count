resource "aws_dynamodb_table" "counter_table" {
  name = "counters"
  read_capacity = 10
  write_capacity = 10
  hash_key = "counter_id"

  attribute {
    name = "counter_id"
    type = "S"
  }
}
