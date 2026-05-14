package api

import (
	"net/http"
	"path/filepath"
)

func (rt *_router) Handler() http.Handler {
	// Serve media
	rt.router.ServeFiles(filepath.Join(rt.media, "*filepath"), http.Dir(rt.rootMedia))

	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)

	// Groups.
	// TODO: Implement removePhoto (which sets the default one).
	rt.router.POST("/groups", rt.authorize(rt.createGroup))
	rt.router.POST("/groups/:groupId", rt.authorize(rt.addToGroup))
	rt.router.GET("/groups/:groupId", rt.authorize(rt.getGroup))
	rt.router.DELETE("/groups/:groupId", rt.authorize(rt.deleteGroup))
	rt.router.PUT("/groups/:groupId/name", rt.authorize(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.authorize(rt.setGroupPhoto))
	rt.router.POST("/groups/:groupId/message", rt.authorize(rt.sendMessageToGroup))
	rt.router.POST("/groups/:groupId/fmessage", rt.authorize(rt.forwardMessageToGroup))
	rt.router.DELETE("/groups/:groupId/user", rt.authorize(rt.leaveGroup))
	rt.router.DELETE("/groups/:groupId/user/:userId", rt.authorize(rt.removeUser))

	// Conversations.
	rt.router.GET("/conversations", rt.authorize(rt.getMyConversations))
	rt.router.GET("/conversations/:conversationId", rt.authorize(rt.getConversation))
	rt.router.GET("/conversations/:conversationId/members", rt.authorize(rt.getMembers))

	// Reactions.
	rt.router.POST("/reactions/:messageId", rt.authorize(rt.addReaction))
	rt.router.DELETE("/reactions/:messageId", rt.authorize(rt.deleteReaction))
	rt.router.GET("/reactions/:messageId", rt.authorize(rt.getReactions))

	// Comments.
	rt.router.DELETE("/comments/:messageId", rt.authorize(rt.uncommentMessage))
	rt.router.POST("/comments/:messageId", rt.authorize(rt.commentMessage))

	// Attachments.
	rt.router.GET("/attachments/:messageId", rt.authorize(rt.getAttachment))

	// Messages.
	rt.router.PUT("/messages/:messageId/status", rt.authorize(rt.updateStatus))
	rt.router.GET("/messages/:messageId/status", rt.authorize(rt.getStatus))
	rt.router.GET("/messages/:messageId", rt.authorize(rt.getMessage))
	rt.router.DELETE("/messages/:messageId", rt.authorize(rt.deleteMessage))

	// Users.
	// TODO: Implement removePhoto (which sets the default one).
	rt.router.GET("/users", rt.authorize(rt.getUsers))
	rt.router.GET("/users/:userId", rt.authorize(rt.getUserById))
	rt.router.DELETE("/users/:userId", rt.authorize(rt.deleteUser))
	rt.router.PUT("/users/:userId/name", rt.authorize(rt.setMyUserName))
	rt.router.PUT("/users/:userId/photo", rt.authorize(rt.setMyPhoto))
	rt.router.POST("/users/:userId/message", rt.authorize(rt.sendMessage))
	rt.router.POST("/users/:userId/fmessage", rt.authorize(rt.forwardMessage))

	// Authentication.
	rt.router.POST("/session", rt.doLogin)

	return rt.router
}
