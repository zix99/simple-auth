<template>
  <Card title="Create Account">
    <form action="/signin" method="POST">

      <div class="field">
        <label class="label">Email</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input is-danger" type="email" placeholder="Email input" />
          <span class="icon is-small is-left">
            <i class="fas fa-envelope" />
          </span>
          <span class="icon is-small is-right">
            <i class="fas fa-exclamation-triangle" />
          </span>
        </div>
        <p class="help is-danger">Email address invalid</p>
      </div>

      <div class="field">
        <label class="label">Username</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="text" placeholder="Text Input">
          <span class="icon is-small is-left">
            <i class="fas fa-user" />
          </span>
        </div>
      </div>

      <div class="field">
        <label class="label">Password</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="password" placeholder="Password" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
      </div>
      <div class="field">
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="password" placeholder="Re-Enter Password" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
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
import Card from './components/card.vue';

export default {
  components: {
    Card,
  },
};
</script>
