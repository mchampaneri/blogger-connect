<template>  
    <div>    
        <BlockUI message="Saving to blogger" v-if="block"></BlockUI>
        <div class="form-group">
          <label class="form-label" for="input-example-1">Title of the post</label>
          <input class="form-input" type="text" v-model="title">
        </div>
        <wysiwyg v-model="content" />  
        <div class="form-group">
            <button class="btn btn-default" @click="save()">Save</button>
            <button class="btn btn-danger" @click="state()">{{nextStatus}}</button>
        </div>       
   </div>
</template>

<script>
export default {
 
    props : ["data","name", "blogid","postid", "status"],

    computed: {
        nextStatus: function () {
        return this.statusTouse == "DRAFT" ? "LIVE": "DRAFT"
        }
    },

    mounted(){
        console.log("desk has been mounted ")
        if (this.data != ""){
            this.content = this.data
        }
        if (this.name != ""){
            this.title = this.name
        }
        if (this.status != ""){
            this.statusTouse = this.status
        }
    },

    data(){
        return{
            content:"",
            title:"",
            block:false,
            statusTouse:"",
        }
    },

    methods:{

        state: function(){
             this.block = true
              axios.post("/chage-state/post",{
                blogid:this.blogid,
                postid:this.postid,
                status:this.statusTouse,
              }).then((succ)=> {
                this.statusTouse = succ.data.status;              
                this.block = false
                console.log("Opertion done succefully :"+ succ)
            }).catch((err)=>{
                 this.block = false
                console.log("Operation failed :"+ err)
            })
        },

        save: function(){
            this.block = true
            axios.post("/save/post",{
                blogid:this.blogid,
                postid:this.postid,
                title:this.title,
                content:this.content,
                status:this.status,
            }).then((succ)=> {
                this.content = succ.data.content;
                this.title = succ.data.title;
                if (this.postid == ""){
                      window.location.href ="/explore/blog/"+this.blogid+"/post/"+succ.data.postid
                }
                this.block = false
                console.log("Opertion done succefully :"+ succ)
            }).catch((err)=>{
                this.block = false
                console.log("Operation failed :"+ err)
            })
        }
    },

}
</script>

