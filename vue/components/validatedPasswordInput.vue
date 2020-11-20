<template>
  <div class="mb-3">
    <div class="field">
      <label class="label">{{title}}</label>
      <div class="control has-icons-left has-icons-right">
        <input class="input" type="password" placeholder="Password" v-model="password1" />
        <span class="icon is-small is-left">
          <fa-icon icon="lock" />
        </span>
      </div>
      <p class="help is-danger" v-if="!validPassword">
        Expected password to be between {{minlength}} and {{maxlength}} long,
        and a score greater than {{minstrength}} (Currently {{strength.score || 0}})
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
    validPassword() {
      return this.password1.length >= this.minlength
        && this.password1.length <= this.maxlength
        && this.strength.score >= this.minstrength;
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
  },
};
</script>
