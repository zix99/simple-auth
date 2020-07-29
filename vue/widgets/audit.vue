<template>
  <div>
    <LoadingBanner :promise="loadingPromise">Fetching Audit...</LoadingBanner>
    <table class="table" v-if="records">
      <thead>
        <tr>
          <th>Date</th>
          <th>Level</th>
          <th>Module</th>
          <th>Message</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="record in this.records" :key="record.ts">
          <td>{{formatDate(record.ts)}}</td>
          <td>{{record.level}}</td>
          <td>{{record.module}}</td>
          <td>{{record.message}}</td>
        </tr>
      </tbody>
    </table>
    <div class="has-text-right">
      <button class="button is-light" @click="prevPage">Previous</button>
      <button class="button is-light" @click="nextPage">Next</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import dayjs from 'dayjs';
import LoadingBanner from '../components/loadingBanner.vue';

export default {
  components: {
    LoadingBanner,
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
  methods: {
    formatDate(s) {
      return dayjs(s).format('lll');
    },
    nextPage() {
      this.offset += this.limit;
    },
    prevPage() {
      this.offset -= this.limit;
      if (this.offset < 0) this.offset = 0;
    },
    fetchData() {
      this.loadingPromise = axios.get('/api/ui/account/audit', { params: { offset: this.offset, limit: this.limit } })
        .then((resp) => {
          this.records = resp.data.records;
        });
    },
  },
};
</script>
