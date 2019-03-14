const multiplyByTwo = async number => number * 2

const main = async () => {
  console.log(await multiplyByTwo(1))
  console.log(await multiplyByTwo(2))
  console.log(await multiplyByTwo(3))
}

main()
