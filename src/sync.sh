rsync -avp --exclude=.git  --exclude=.idea --exclude=vendor --exclude=.DS_Store --exclude='/service_*' --exclude='/tool_*'  --exclude='/task_*' . media1:~/workspace/rds_slowq_query/src
