package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)

	// groups
	rt.router.POST("/groups", rt.validateAuthorization(rt.createGroup))
	rt.router.GET("/groups/:groupId", rt.validateAuthorization(rt.getGroup))
	rt.router.DELETE("/groups/:groupId", rt.validateAuthorization(rt.deleteGroup))
	rt.router.POST("/groups/:groupId", rt.validateAuthorization(rt.addToGroup))
	rt.router.PUT("/groups/:groupId/name", rt.validateAuthorization(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.validateAuthorization(rt.setGroupPhoto))
	rt.router.POST("/groups/:groupId/message", rt.validateAuthorization(rt.sendMessageToGroup))
	rt.router.POST("/groups/:groupId/fmessage", rt.validateAuthorization(rt.forwardMessageToGroup))
	rt.router.DELETE("/groups/:groupId/user", rt.validateAuthorization(rt.leaveGroup))
	rt.router.DELETE("/groups/:groupId/user/:userId", rt.validateAuthorization(rt.removeUser))

	// conversations
	rt.router.GET("/conversations", rt.validateAuthorization(rt.getMyConversations))
	rt.router.GET("/conversations/:conversationId", rt.validateAuthorization(rt.getConversation))
	rt.router.GET("/conversations/:conversationId/members", rt.validateAuthorization(rt.getMembers))

	// reactions
	rt.router.POST("/reactions/:messageId", rt.validateAuthorization(rt.addReaction))
	rt.router.GET("/reactions/:messageId", rt.validateAuthorization(rt.getReactions))
	rt.router.DELETE("/reactions/:messageId", rt.validateAuthorization(rt.deleteReaction))

	// comments
	rt.router.DELETE("/comments/:messageId", rt.validateAuthorization(rt.deleteMessage))
	rt.router.POST("/comments/:messageId", rt.validateAuthorization(rt.commentMessage))

	// messages
	rt.router.GET("/messages/:messageId", rt.validateAuthorization(rt.getMessage))
	rt.router.DELETE("/messages/:messageId", rt.validateAuthorization(rt.deleteMessage))
	rt.router.GET("/messages/:messageId/status", rt.validateAuthorization(rt.getStatus))

	// user
	rt.router.POST("/user/:userId/message", rt.validateAuthorization(rt.sendMessage))
	rt.router.POST("/user/:userId/fmessage", rt.validateAuthorization(rt.forwardMessage))
	rt.router.GET("/user/:userId", rt.validateAuthorization(rt.getUserById))
	rt.router.PUT("/user/:userId/name", rt.validateAuthorization(rt.setMyUserName))
	rt.router.PUT("/user/:userId/photo", rt.validateAuthorization(rt.setMyPhoto))
	rt.router.GET("/user", rt.validateAuthorization(rt.getUsers))
	rt.router.DELETE("/user", rt.validateAuthorization(rt.deleteUser))

	rt.router.POST("/session", rt.doLogin)

	return rt.router
}
