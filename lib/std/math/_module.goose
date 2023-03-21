native fn sin(x)
native fn cos(x)
native fn tan(x)
native fn asin(x)
native fn acos(x)
native fn atan(x)
native fn atan2(y, x)
native fn sinh(x)
native fn cosh(x)
native fn tanh(x)

native fn round(x)
native fn floor(x)
native fn ceil(x)
native fn sqrt(x)

native fn log(x)
native fn log2(x)
native fn log10(x)
native fn exp(x)
native fn pow(x, y)

native fn abs(x)
native fn sign(x)
native fn min(x, y)
native fn max(x, y)
native fn clamp(x, min, max)

native fn parseInt(x, base)
// native fn parseFloat(x)

native const pi
native const e
native const phi
native const sqrt2
native const ln2
native const log2e
native const ln10
native const log10e

export {
  sin, cos, tan, asin, acos, atan, atan2, sinh, cosh, tanh,
  round, floor, ceil, sqrt,
  log, log2, log10, exp, pow,
  abs, sign, min, max, clamp,

  parseInt,

  pi, e, phi, sqrt2, ln2, log2e, ln10, log10e
}
