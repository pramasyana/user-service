# USER SERVICE PROTOCOLS

  Contains all USER SERVICE API Protocols

# Protocol Buffer (Protobuf)
  Protocol buffers are Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data – think XML, but smaller, faster, and simpler.

  Requirements

   - Golang version 1.7+
   - https://github.com/google/protobuf/releases

  How to generate Go code from protobuf file:

  - Generate code from all protobuf file:
    ```shell
    $ make all
    ```

  - Generate Health Code:
    ```shell
    $ make health
    ```

  results will be in the proto-go folder github.com/Bhinneka/user-service/proto-go/
