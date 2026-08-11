(function() {
  var chbCheckin = {
    init: function() {
      this.addCheckinButton();
      this.loadStatus();
    },

    addCheckinButton: function() {
      var container = document.querySelector('.user-menu .menu-panel');
      if (!container) {
        setTimeout(function() { chbCheckin.addCheckinButton(); }, 1000);
        return;
      }

      var btn = document.createElement('button');
      btn.className = 'btn btn-default chb-checkin-btn';
      btn.innerHTML = '每日签到';
      btn.title = '签到获取 CHB 积分';
      btn.addEventListener('click', function() { chbCheckin.doCheckin(this); });
      container.insertBefore(btn, container.firstChild);
    },

    loadStatus: function() {
      var xhr = new XMLHttpRequest();
      xhr.open('GET', '/chb/checkin/status', true);
      xhr.setRequestHeader('Content-Type', 'application/json');
      xhr.onload = function() {
        if (xhr.status === 200) {
          var resp = JSON.parse(xhr.responseText);
          if (resp.data && resp.data.checked_in_today) {
            chbCheckin.updateButtonStatus(true);
          }
        }
      };
      xhr.send();
    },

    doCheckin: function(btn) {
      btn.disabled = true;
      btn.innerHTML = '签到中...';

      var xhr = new XMLHttpRequest();
      xhr.open('POST', '/chb/checkin', true);
      xhr.setRequestHeader('Content-Type', 'application/json');
      xhr.onload = function() {
        var resp = JSON.parse(xhr.responseText);
        if (resp.code === 0 && resp.data && resp.data.status === 'completed') {
          btn.innerHTML = '已签到 +' + (resp.data.final_amount || 0) + ' CHB';
          btn.className = 'btn btn-default chb-checkin-btn checked-in';
          btn.disabled = true;
        } else if (resp.data && resp.data.status === 'rejected') {
          if (resp.data.reject_reason === 'already_checked_in') {
            btn.innerHTML = '今日已签到';
            btn.className = 'btn btn-default chb-checkin-btn checked-in';
          } else {
            btn.innerHTML = '签到失败';
            btn.className = 'btn btn-default chb-checkin-btn failed';
            btn.disabled = false;
          }
        } else {
          btn.innerHTML = '签到失败';
          btn.className = 'btn btn-default chb-checkin-btn failed';
          btn.disabled = false;
        }
      };
      xhr.onerror = function() {
        btn.innerHTML = '网络错误';
        btn.className = 'btn btn-default chb-checkin-btn failed';
        btn.disabled = false;
      };
      xhr.send(JSON.stringify({}));
    },

    updateButtonStatus: function(checkedIn) {
      var btn = document.querySelector('.chb-checkin-btn');
      if (btn) {
        if (checkedIn) {
          btn.innerHTML = '今日已签到';
          btn.className = 'btn btn-default chb-checkin-btn checked-in';
          btn.disabled = true;
        }
      }
    }
  };

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function() { chbCheckin.init(); });
  } else {
    chbCheckin.init();
  }
})();
