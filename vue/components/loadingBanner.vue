<template>
  <div>
    <article class="message is-danger" v-if="state === States.REJECTED && error">
      <div class="message-body">
        <fa-icon icon="exclamation-triangle" /> <slot name="error" :error="error">{{error}}</slot>
      </div>
    </article>
    <article class="message is-info" v-if="state === States.LOADING">
      <div class="message-body">
        <fa-icon icon="circle-notch" spin /> <slot>Loading...</slot>
      </div>
    </article>
    <article class="message is-success" v-if="state === States.RESOLVED && this.$slots.success">
      <div class="message-body">
        <fa-icon icon="check" /> <slot name="success"></slot>
      </div>
    </article>
  </div>
</template>

<script>

const States = Object.freeze({
  NONE: 0,
  LOADING: 1,
  RESOLVED: 2,
  REJECTED: 3,
});

export default {
  props: {
    promise: null,
    codes: { default: () => ({}) },
  },
  data() {
    return {
      States,
      state: States.NONE,
      error: null,
    };
  },
  created() {
    this.updatePromise();
  },
  watch: {
    promise() {
      this.updatePromise();
    },
    state() {
      this.$emit('state', this.state);
    },
  },
  methods: {
    updatePromise() {
      if (!this.promise) {
        this.state = States.NONE;
      } else {
        this.error = null;
        this.promise.pending = true;

        // Delay showing the loader so it's less likely to blink in/out so quickly
        const timer = setTimeout(() => {
          this.state = States.LOADING;
        }, 100);

        this.promise
          .then(() => {
            this.state = States.RESOLVED;
          })
          .catch((err) => {
            this.error = this.extractErrorMessage(err);
            this.state = States.REJECTED;
          })
          .then(() => {
            this.promise.pending = false;
            clearTimeout(timer);
          });
      }
    },
    extractErrorMessage(err) {
      if (err.response && err.response.data && err.response.data.error === true) {
        const errdata = err.response.data;
        if (errdata.reason && this.codes[errdata.reason]) {
          return this.codes[errdata.reason];
        }
        return errdata.message;
      }
      return err.message;
    },
  },
};
</script>
