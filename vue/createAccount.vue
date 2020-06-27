<template>
  <Card title="Create Account">
    <form action="/signin" method="POST">

      <div class="field">
        <label class="label">Email</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" :class="{ 'is-danger': !validEmail }" type="email" placeholder="Email input" v-model="email" />
          <span class="icon is-small is-left">
            <i class="fas fa-envelope" />
          </span>
          <span class="icon is-small is-right" v-if="!validEmail">
            <i class="fas fa-exclamation-triangle" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validEmail">Email address invalid</p>
      </div>

      <div class="field">
        <label class="label">Username</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="text" placeholder="Text Input" v-model="username">
          <span class="icon is-small is-left">
            <i class="fas fa-user" />
          </span>
          <span class="icon is-small is-right" v-if="!validUsername">
            <i class="fas fa-exclamation-triangle" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validUsername">Expected username to be between {{usernameMinLength}} and {{usernameMaxLength}} long</p>
      </div>

      <div class="field">
        <label class="label">Password</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="password" placeholder="Password" v-model="password1" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validPassword">
          Expected password to be between {{passwordMinLength}} and {{passwordMaxLength}} long, and a score greater than 2 (Currently {{strength.score}})
        </p>
        <p class="help" v-if="strength.feedback">
          {{strength.feedback.warning}}
          <ul>
            <li v-for="suggestion in strength.feedback.suggestions" :key="suggestion">{{suggestion}}</li>
          </ul>
        </p>
      </div>
      <div class="field">
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="password" placeholder="Re-Enter Password" v-model="password2" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
        <p class="help is-danger" v-if="!passwordMatch">Password does not match</p>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link">Submit</button>
        </div>
      </div>
    </form>
  </Card>
</template>

<script>
import validator from 'validator';
import zxcvbn from 'zxcvbn';
import Card from './components/card.vue';

export default {
  props: {
    usernameMinLength: { type: Number, default: 1 },
    usernameMaxLength: { type: Number, default: 999 },
    passwordMinLength: { type: Number, default: 1 },
    passwordMaxLength: { type: Number, default: 999 },
  },
  data() {
    return {
      // input
      email: '',
      username: '',
      password1: '',
      password2: '',

      // responsive
      strength: {},
    };
  },
  components: {
    Card,
  },
  computed: {
    validEmail() {
      return validator.isEmail(this.email);
    },
    validUsername() {
      return this.username.length >= this.usernameMinLength && this.username.length <= this.usernameMaxLength;
    },
    passwordMatch() {
      return this.password1 === this.password2;
    },
    validPassword() {
      return this.password1.length >= this.passwordMinLength
        && this.password1.length <= this.passwordMaxLength
        && this.strength.score >= 2;
    },
  },
  watch: {
    password1() {
      this.strength = zxcvbn(this.password1);
    },
  },
};
</script>
