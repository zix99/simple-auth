<template>
  <div>
    <LoadingBanner :promise="loadingPromise" :codes="errorCodes">Updating password...</LoadingBanner>
    <Message v-if="success" type="is-success">Password updated!</Message>
    <div v-if="!success">
      <div class="field" v-if="requireOldPassword">
        <label class="label">Old Password</label>
        <div class="control has-icons-left">
          <input class="input" type="password" placeholder="Password" v-model="oldpassword" v-focus />
          <span class="icon is-small is-left">
            <fa-icon icon="lock" />
          </span>
        </div>
      </div>

      <ValidatedPasswordInput title="New Password" v-model="password" @valid="validPassword = $event" @enter="submitClick" />
      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="!validInput">Update</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import ValidatedPasswordInput from '../components/validatedPasswordInput.vue';
import LoadingBanner from '../components/loadingBanner.vue';
import Message from '../components/message.vue';

export default {
  components: {
    ValidatedPasswordInput,
    LoadingBanner,
    Message,
  },
  data() {
    return {
      loadingPromise: null,
      success: false,
      requireOldPassword: true,

      validPassword: false,
      password: '',
      oldpassword: '',

      errorCodes: {
        'invalid-credentials': 'Your old password is invalid',
        'unsatisfied-stipulations': 'Your account has an unsatisfied stipulation on it',
      },
    };
  },
  mounted() {
    axios.get('api/ui/account/password')
      .then((resp) => {
        this.requireOldPassword = resp.data.requireOldPassword;
      });
  },
  methods: {
    submitClick() {
      const data = {
        oldpassword: this.oldpassword,
        newpassword: this.password,
      };
      this.loadingPromise = axios.post('api/ui/account/password', data)
        .then(() => {
          this.success = true;
          setTimeout(() => this.$emit('submitted'), 1500);
        });
    },
  },
  computed: {
    validInput() {
      return this.validPassword;
    },
  },
};
</script>
