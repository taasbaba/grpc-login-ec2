<template>
    <div>
        <h1>Hi {{account.user.username}}!</h1>
        <p>You're logged in with Gaming Heaven!! <router-link to="/login">Logout</router-link></p>
    </div>
</template>

<script>
import { mapState, mapActions } from 'vuex'

export default {
    computed: {
        ...mapState({
            account: state => state.account,
            users: state => state.users.all
        })
    },
    created () {
        this.getAllUsers();
        console.log('In homePage::Acount', this.account)
    },
    methods: {
        ...mapActions('users', {
            getAllUsers: 'getAll',
            deleteUser: 'delete'
        })
    }
};
</script>