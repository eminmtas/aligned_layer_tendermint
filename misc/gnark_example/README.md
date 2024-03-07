# Using Groth16 with gnark

This example demonstrates how to use gnark.go to run a [Groth16](https://www.zeroknowledgeblog.com/index.php/groth16) proof. It follows the documentation to instantiate the circuit and verify it.

If you only want to run the verification step, you'll need the proof, verifying key, and public witness. You may need to serialize these components beforehand.

## How to Run It

```sh
go run gnark_example.go
```
