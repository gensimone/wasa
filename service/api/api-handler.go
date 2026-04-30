package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)


	// Groups.
	rt.router.GET("/groups/:groupId", rt.validateAuthorization(rt.getGroup))
	rt.router.DELETE("/groups/:groupId", rt.validateAuthorization(rt.deleteGroup))
	rt.router.PUT("/name/group/:groupId", rt.validateAuthorization(rt.setGroupName))
	rt.router.PUT("/photo/group/:groupId", rt.validateAuthorization(rt.setGroupPhoto))
	rt.router.DELETE("/leave/:groupId", rt.validateAuthorization(rt.leaveGroup))
	// TODO:
	rt.router.POST("/groups/:groupId", rt.validateAuthorization(rt.addToGroup))
	rt.router.POST("/groups", rt.validateAuthorization(rt.createGroup))
	rt.router.POST("/forward/groups/:groupId", rt.validateAuthorization(rt.forwardMessageToGroup))
	rt.router.POST("/send/groups/:groupId", rt.validateAuthorization(rt.sendMessageToGroup))
	rt.router.DELETE("/remove/:groupId/user/:userId", rt.validateAuthorization(rt.removeUser))

	// Conversations.
	rt.router.GET("/members/:conversationId", rt.validateAuthorization(rt.getMembers))
	rt.router.GET("/conversations", rt.validateAuthorization(rt.getMyConversations))
	rt.router.GET("/conversations/:conversationId", rt.validateAuthorization(rt.getConversation))

	// Reactions.
	// TODO:
	rt.router.POST("/reactions/:messageId", rt.validateAuthorization(rt.addReaction))
	rt.router.GET("/reactions/:messageId", rt.validateAuthorization(rt.getReactions))
	rt.router.DELETE("/reactions/:messageId", rt.validateAuthorization(rt.deleteReaction))

	// Messages.
	rt.router.GET("/messages/:messageId", rt.validateAuthorization(rt.getMessage))
	rt.router.DELETE("/messages/:messageId", rt.validateAuthorization(rt.deleteMessage))
	rt.router.DELETE("/comments/:messageId", rt.validateAuthorization(rt.deleteMessage))
	rt.router.POST("/send/user/:userId", rt.validateAuthorization(rt.sendMessage))
	rt.router.POST("/forward/user/:userId", rt.validateAuthorization(rt.forwardMessage))
	rt.router.POST("/comments/:messageId", rt.validateAuthorization(rt.commentMessage))
	rt.router.GET("/status/:messageId", rt.validateAuthorization(rt.getStatus))

	// Users.
	rt.router.GET("/user/:userId", rt.validateAuthorization(rt.getUserById))
	rt.router.GET("/user", rt.validateAuthorization(rt.getUsers))
	rt.router.DELETE("/user", rt.validateAuthorization(rt.deleteUser))
	rt.router.PUT("/name/user", rt.validateAuthorization(rt.setMyUserName))
	rt.router.PUT("/photo/user", rt.validateAuthorization(rt.setMyPhoto))
	rt.router.POST("/session", rt.doLogin)

	return rt.router
}
