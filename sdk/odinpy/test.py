from odinpy import Client
c = Client()
s = c.new_sandbox()

s.upsert_directory("./dir")

f = s.get_file("./dir/file.txt")

f.upsert(patch="--- hello.c\t2025-03-16 20:34:09.998561417 +0530\n+++ hello_new.c\t2025-03-16 20:32:46.721067247 +0530\n@@ -0,0 +1,6 @@\n+#include <stdio.h>\n+\n+int main(int argc, char *argv[]) {\n+    printf(\"Hello World\\n\");\n+    return 0;\n+}\n\\ No newline at end of file\n")


content = f.read().content
print(content)
print(s.read_directory("./dir"))
s.delete_file("./dir/file.txt")
content = f.read().content
print(content or "File does not exist")
s.delete_directory("./dir")
print(s.read_directory("./dir"))