<template>  
    <div>
        <div class="form-group">
          <label class="form-label" for="input-example-1">Name</label>
          <input class="form-input" type="text" v-model="title">
        </div>
        <wysiwyg v-model="content" /> 
        <div class="form-group">
            <button class="btn btn-default" @click="save()">Save</button>
        </div>
   </div>
</template>

<script>
export default {
 

    props : ["data","name", "blogid","postid"],

    mounted(){
        console.log("desk has been mounted ")
        if (this.data != ""){
            this.content = this.data
        }
        if (this.name != ""){
            this.title = this.name
        }
    },

    data(){
        return{
            content:"",
            title:""
        }
    },

    methods:{
        save: function(){
            axios.post("/save/post",{
                blogid:this.blogid,
                postid:this.postid,
                title:this.title,
                content:this.content
            }).then((succ)=> {
                this.content = succ.data.content;
                this.title = succ.data.title;
                console.log("Opertion done succefully :"+ succ)
            }).catch((err)=>{
                console.log("Operation failed :"+ err)
            })
        }
    }

}
</script>

