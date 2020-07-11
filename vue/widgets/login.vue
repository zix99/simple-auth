<template>
  <div>
    <article class="message is-danger" v-if="error">
      <div class="message-body">{{error}}</div>
    </article>

    <div v-if="state === 'login'">
      <div class="field">
        <label class="label">Username</label>
        <div class="control has-icons-left">
          <input class="input" type="text" placeholder="Text Input" v-model="username">
          <span class="icon is-small is-left">
            <i class="fas fa-user" />
          </span>
        </div>
      </div>

      <div class="field">
        <label class="label">Password</label>
        <div class="control has-icons-left">
          <input class="input" type="password" placeholder="Password" v-model="password" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick">Login</button>
        </div>
      </div>
      <a href="#" @click.prevent="forgotPassword">Forgot Password?</a>
    </div>

    <div v-if="state === 'forgot'">
      Forgot pass
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      error: null,
      username: '',
      password: '',
      state: 'login',
    };
  },
  methods: {
    submitClick() {
      const postData = {
        username: this.username,
        password: this.password,
      };
      axios.post('api/ui/login', postData)
        .then(() => {
          this.$emit('loggedIn');
        }).catch((err) => {
          this.error = err.message;
        });
    },
    forgotPassword() {
      this.$emit('forgotPassword');
    },
  },
};
</script>
