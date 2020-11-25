export default function debounce(fn, delay) {
  let timer = 0;
  return function debounced(...args) {
    const context = this;
    clearTimeout(timer);
    timer = setTimeout(() => fn.apply(context, args), delay);
  };
}
