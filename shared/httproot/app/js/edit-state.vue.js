export const EditState = {

    data() {
        return {
            // state/datet at the last load/save, null while nothing is loaded
            savedState: null
        }
    },

    methods: {

        getState() {
            console.log("Not implemented, please override getState()");
            return null;
        },

        saveState() {
            const state = this.getState();
            this.savedState = state != null ? JSON.stringify(state) : null;
        },

        stateHasChanged() {
            const state = this.getState();
            const check = state != null ? JSON.stringify(state) : null;
            return check !== this.savedState;
        },

        async confirmLeave() {
            if(!this.stateHasChanged())
                return true;
            return confirm("Leave without saving?");
        }
    }

}