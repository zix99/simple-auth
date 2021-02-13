<template>
  <div>
    <LoadingBanner :promise="loadingPromise">Fetching Audit...</LoadingBanner>
    <table class="table is-narrow is-hoverable is-fullwidth" v-if="records">
      <thead>
        <tr>
          <th>Date</th>
          <th class="is-hidden-mobile">Level</th>
          <th class="is-hidden-mobile">Module</th>
          <th>Message</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="record in this.records" :key="record.ts">
          <td><ShortDate :date="record.ts" /></td>
          <td :class="levelToClass(record.level)" class="is-hidden-mobile">{{record.level}}</td>
          <td class="is-hidden-mobile">{{record.module}}</td>
          <td>{{record.message}}</td>
        </tr>
      </tbody>
    </table>
    <div class="has-text-right">
      <button class="button is-light" @click="prevPage">Previous</button>
      <button class="button is-light" @click="nextPage" :disabled="!hasNextButton">Next</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import LoadingBanner from '../components/loadingBanner.vue';
import ShortDate from '../components/shortdate.vue';

export default {
  components: {
    LoadingBanner,
    ShortDate,
  },
  data() {
    return {
      offset: 0,
      limit: 10,
      records: null,
      loadingPromise: null,
    };
  },
  created() {
    this.fetchData();
  },
  watch: {
    offset() {
      this.fetchData();
    },
    limit() {
      this.fetchData();
    },
  },
  computed: {
    hasNextButton() {
      return this.records.length >= this.limit;
    },
  },
  methods: {
    nextPage() {
      this.offset += this.limit;
    },
    prevPage() {
      this.offset -= this.limit;
      if (this.offset < 0) this.offset = 0;
    },
    fetchData() {
      this.loadingPromise = axios.get('api/v1/account/audit', { params: { offset: this.offset, limit: this.limit } })
        .then((resp) => {
          this.records = resp.data.records;
        });
    },
    levelToClass(lvl) {
      const level = lvl.toLowerCase();
      if (level === 'alert') {
        return 'is-danger';
      }
      if (level === 'warn') {
        return 'is-warning';
      }
      return 'is-info';
    },
  },
};
</script>
