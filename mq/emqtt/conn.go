package emqtt

type Emqtt struct {
	Host      string `json:"host" description:"EMQTT地址"`
	Port      int    `json:"port" description:"EMQTT端口"`
	Username  string `json:"username" description:"EMQTT访问用户名"`
	Password  string `json:"password" description:"EMQTT访问密码"`
	TopicName string `json:"topicName" description:"EMQTT TopicName"`
}
