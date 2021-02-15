<template>
  <div class="mb-3">
    <div class="field">
      <label class="label">{{title}}</label>
      <div class="control has-icons-left has-icons-right">
        <input class="input" :class="strengthClass" type="password" placeholder="Password" v-model="password1" @keypress.enter="enterPress" />
        <span class="icon is-small is-left">
          <fa-icon icon="lock" />
        </span>
      </div>
      <p class="help is-danger" v-if="!validPasswordLength">
        Expected password to be between {{minlength}} and {{maxlength}} long.
      </p>
      <p class="help is-danger" v-if="!validPasswordStrength && validPasswordLength">
        Password not strong enough.
      </p>
      <p class="help" v-if="strength.feedback && validPasswordLength">
        {{strength.feedback.warning}}
        <ul>
          <li v-for="suggestion in strength.feedback.suggestions" :key="suggestion">{{suggestion}}</li>
        </ul>
      </p>
    </div>
    <div class="field">
      <div class="control has-icons-left has-icons-right">
        <input class="input" :class="{ 'is-success': passwordMatch }" type="password" placeholder="Re-Enter Password" v-model="password2" @keypress.enter="enterPress" />
        <span class="icon is-small is-left">
          <fa-icon icon="lock" />
        </span>
      </div>
      <p class="help is-danger" v-if="!passwordMatch">Password does not match</p>
    </div>
  </div>
</template>

<script>
import zxcvbn from 'zxcvbn';

export default {
  props: {
    minlength: { default: 1 },
    maxlength: { default: 999 },
    minstrength: { default: 2 },
    value: null,
    title: { default: 'Password' },
  },
  data() {
    return {
      password1: '',
      password2: '',
      strength: { score: 0 },
    };
  },
  computed: {
    passwordMatch() {
      return this.password1 === this.password2;
    },
    validPasswordLength() {
      return this.password1.length >= this.minlength
        && this.password1.length <= this.maxlength;
    },
    validPasswordStrength() {
      return this.strength.score >= this.minstrength;
    },
    validPassword() {
      return this.validPasswordLength && this.validPasswordStrength;
    },
    strengthClass() {
      if (this.strength.score < this.minstrength) return 'is-danger';
      if (this.strength.score === 4) return 'is-success';
      if (this.strength.score === 3) return 'is-warning';
      return 'is-warning';
    },
  },
  watch: {
    password1() {
      this.strength = zxcvbn(this.password1);
      this.onChange();
    },
    password2() {
      this.onChange();
    },
  },
  methods: {
    onChange() {
      this.$emit('valid', this.passwordMatch && this.validPassword);
      this.$emit('input', this.password1);
    },
    enterPress() {
      if (this.validPassword && this.passwordMatch) {
        this.$emit('enter');
      }
    },
  },
};
</script>
