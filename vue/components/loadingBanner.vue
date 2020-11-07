<template>
  <div>
    <article class="message is-danger" v-if="error">
      <div class="message-body">
        <i class="fas fa-exclamation-triangle"></i> {{error}}
      </div>
    </article>
    <article class="message is-info" v-if="loading">
      <div class="message-body">
        <i class="fas fa-circle-notch fa-spin"></i> <slot />
      </div>
    </article>
  </div>
</template>

<script>
export default {
  props: {
    promise: null,
    codes: { default: () => ({}) },
  },
  data() {
    return {
      loading: false,
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
  },
  methods: {
    updatePromise() {
      if (!this.promise) {
        this.loading = false;
      } else {
        this.error = null;
        this.promise.pending = true;

        // Delay showing the loader so it's less likely to blink in/out so quickly
        const timer = setTimeout(() => {
          this.loading = true;
        }, 100);
        this.promise
          .catch((err) => {
            this.error = this.extractErrorMessage(err);
          })
          .then(() => {
            this.promise.pending = false;
            this.loading = false;
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
