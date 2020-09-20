## jsonlogic-go

A [jsonlogic](http://jsonlogic.com) library in golang. 

It is meant to be a 'stricter' version comparing to the [js](https://github.com/jwadhams/json-logic-js/) version:
It should behave the same as the js version, if not then an error should be returned.

This means that if an expression is evaluated successful in server side using this library, 
then it is expected to be evaluated to the same result in client side. But the reverse direction maybe not true.

### Notable restrictions

- `==`/`!=` operations are not supported. Use strict version instead: `===`/`!==`
- Many operations will check the minimal number of params. See doc for detail. Some examples:
  - `{"var":[]}` is ok in js. But not ok in jsonlogic-go. (`{"var":null}` or `{"var":[null]}` is ok though)
  - `{"===":["x"]}` is ok in js. But jsonlogic-go requires at least two params.
- Many operations will check the type of params. Some examples:
  - Comparing and equality checking only accepts json primitves (`null/bool/numeric/string`)
    - You can still `{">":[1,"-1"]}` but not `{">":[1,[]]}`
- No `NaN`: `{">":["b",0]}` is evaluated to `false` in js since `"b"` is converted to `NaN`. But here you will get an error since parse error.
  - More detail about comparing in js: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than
