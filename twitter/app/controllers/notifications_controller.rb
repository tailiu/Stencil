class NotificationsController < ApplicationController
    def get
        result = {
            # "params" => params,
            "success" => false,
            "error" => {
            },
        }
        if params[:user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            result["success"] = true
            notifs = Notification.joins("INNER JOIN users ON users.id = notifications.from_user").where(user_id: params[:user_id])
            result["notifs"] = notifs.as_json
            for notif in notifs do
                notif.is_seen = true
                notif.save
            end
        end
        render json: {result: result}
    end
end
