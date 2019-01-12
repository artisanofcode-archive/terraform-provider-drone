resource "drone_secret" "test" {
  repository = "test/test"
  name       = "fancy_pants"
  value      = "correct horse battery staple"
}
