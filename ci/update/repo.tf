resource "drone_repo" "test" {
  repository = "test/test"
  visibility = "public"
  trusted = true
}
